// Package aws2tf ingests JSON from AWS CLI and emits Terraform templates.
// Currently, only security groups are implemented. Poorly, at that.
//
// TODO:
// - learn interfaces
// - make the ipRange.print* methods nicer
package aws2tf

import (
	"fmt"
)

type ipRange struct {
	Description string
	CidrIP      string
}

func (ir ipRange) String() string {
	s := fmt.Sprintf("    cidr_blocks = [\"%s\"]\n", ir.CidrIP)

	if ir.Description != "" {
		s += fmt.Sprintf("    description = \"%s\"\n", ir.Description)
	}

	return s
}

func (ir ipRange) printIngress(ip ipPermissionIngress) string {
	s := "  ingress {\n"
	s += fmt.Sprintf("%s", ip)
	s += fmt.Sprintf("%s", ir)
	s += "  }\n"

	return s
}

func (ir ipRange) printEgress(ip ipPermissionEgress) string {
	s := "  egress {\n"
	s += fmt.Sprintf("%s", ip)
	s += fmt.Sprintf("%s", ir)
	s += "  }\n"

	return s
}

type ipPermission struct {
	FromPort         float64
	ToPort           float64
	IPProtocol       string
	IPRanges         []ipRange
	UserIDGroupPairs []userIDGroupPair
}

func (ip ipPermission) String() string {
	s := fmt.Sprintf("    from_port   = %d\n", int(ip.FromPort))
	s += fmt.Sprintf("    to_port     = %d\n", int(ip.ToPort))
	s += fmt.Sprintf("    protocol    = \"%s\"\n", ip.IPProtocol)

	return s
}

type ipPermissionIngress struct {
	ipPermission
}

type ipPermissionIngresses []ipPermissionIngress

func (ipi ipPermissionIngresses) String() string {
	var s string

	for _, i := range ipi {
		for _, j := range i.IPRanges {
			s += j.printIngress(i)
		}

		for _, k := range i.UserIDGroupPairs {
			s += k.printIngress(i)
		}
	}

	return s
}

type ipPermissionEgress struct {
	ipPermission
}

type ipPermissionEgresses []ipPermissionEgress

func (ipe ipPermissionEgresses) String() string {
	var s string

	for _, i := range ipe {
		for _, j := range i.IPRanges {
			s += j.printEgress(i)
		}

		for _, k := range i.UserIDGroupPairs {
			s += k.printEgress(i)
		}
	}

	return s
}

type userIDGroupPair struct {
	UserID                 string
	GroupID                string
	Description            string
	VpcID                  string
	VpcPeeringConnectionID string
	PeeringStatus          string
}

func (uigp userIDGroupPair) String() string {
	var sg string

	if uigp.VpcID != "" || uigp.VpcPeeringConnectionID != "" || uigp.PeeringStatus != "" {
		sg = uigp.UserID + "/" + uigp.GroupID
	} else {
		sg = uigp.GroupID
	}

	s := fmt.Sprintf("    security_groups = [\"%s\"]\n", sg)

	if uigp.Description != "" {
		s += fmt.Sprintf("    description = \"%s\"\n", uigp.Description)
	}

	return s
}

func (uigp userIDGroupPair) printIngress(ip ipPermissionIngress) string {
	s := "  ingress {\n"
	s += fmt.Sprintf("%s", ip)
	s += fmt.Sprintf("%s", uigp)
	s += "  }\n"

	return s
}

func (uigp userIDGroupPair) printEgress(ip ipPermissionEgress) string {
	s := "  egress {\n"
	s += fmt.Sprintf("%s", ip)
	s += fmt.Sprintf("%s", uigp)
	s += "  }\n"

	return s
}

type tag struct {
	Key   string
	Value string
}

func (t tag) String() string {
	return fmt.Sprintf("    \"%s\" = \"%s\"\n", t.Key, t.Value)
}

type tags []tag

func (t tags) String() string {
	s := "  tags {\n"

	for _, x := range t {
		s += fmt.Sprintf("%s", x)
	}

	s += "  }\n"

	return s
}

type securityGroup struct {
	GroupName           string
	Description         string
	GroupID             string
	VpcID               string
	IPPermissions       ipPermissionIngresses
	IPPermissionsEgress ipPermissionEgresses
	Tags                tags
}

func (sg securityGroup) String() string {
	s := fmt.Sprintf("resource \"aws_security_group\" \"%s-%s\" {\n", sg.VpcID, sg.GroupName)
	s += fmt.Sprintf("  name        = \"%s\"\n", sg.GroupName)
	s += fmt.Sprintf("  description = \"%s\"\n", sg.Description)
	s += fmt.Sprintf("  vpc_id      = \"%s\"\n", sg.VpcID)

	if len(sg.IPPermissions) > 0 {
		s += fmt.Sprintf("%s", sg.IPPermissions)
	}

	if len(sg.IPPermissionsEgress) > 0 {
		s += fmt.Sprintf("%s", sg.IPPermissionsEgress)
	}

	if len(sg.Tags) > 0 {
		s += fmt.Sprintf("%s", sg.Tags)
	}

	s += "}\n"

	return s
}

// SGFile describes a JSON file which would be produced by the AWS CLI, with details of one or more security groups.
// This JSON file would be produced by the 'aws ec2 describe-security-groups' command.
type SGFile struct {
	SecurityGroups []securityGroup
}

func (sgf SGFile) String() string {
	var s string

	for _, sg := range sgf.SecurityGroups {
		s += fmt.Sprintf("%s\n", sg)
	}

	return s
}
