package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func init() {
	os.Setenv("AWS_PROFILE", "celab")

	sess, errNewSess := session.NewSession()
	if errNewSess != nil {
		log.Fatalln("[NewSession()]", errNewSess)
	}

	_, errCredGet := sess.Config.Credentials.Get()
	if errCredGet != nil {
		log.Fatalln("[Config.Credentials.Get()]", errCredGet)
	}

	svcIam = iam.New(sess)
}

var (
	rolesProcessed uint
	sess           *session.Session
	svcIam         *iam.IAM
)

func main() {
	for roleOutput := range genRoles(genRoleNames("/Users/jameslucktaylor/roles.txt")) {
		for ipOutput := range genInstanceProfiles(roleOutput) {
			removeInstanceProfilesFromRole(
				roleIPRemove{
					instanceProfilesOutput: ipOutput,
					roleOutput:             roleOutput,
				},
			)
		}

		for rolePolicies := range genRolePolicies(roleOutput) {
			deletePoliciesFromRole(
				rolePolicyRemove{
					policiesOutput: rolePolicies,
					roleOutput:     roleOutput,
				},
			)
		}

		for attachedPolicies := range genAttachedPolicies(roleOutput) {
			detachPoliciesFromRole(
				rolePolicyDetach{
					attachedPolicies: attachedPolicies,
					roleOutput:       roleOutput,
				},
			)
		}

		deleteRole(aws.StringValue(roleOutput.Role.RoleName))
	}

	fmt.Println("Roles processed:", rolesProcessed)
}

func genRoleNames(filename string) (output chan string) {
	output = make(chan string)

	go func() {
		inFile, _ := os.Open(filename)
		defer inFile.Close()

		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			output <- scanner.Text()
		}

		close(output)
	}()

	return
}

func genRoles(input chan string) (output chan *iam.GetRoleOutput) {
	output = make(chan *iam.GetRoleOutput)

	go func() {
		for roleName := range input {
			rolesProcessed++

			iamRole, errGetRole := svcIam.GetRole(
				&iam.GetRoleInput{
					RoleName: aws.String(roleName),
				},
			)
			if errGetRole != nil {
				if strings.Contains(fmt.Sprint(errGetRole), "NoSuchEntity") {
					log.Printf("[GetRole()] Role doesn't exist: %s\n", roleName)
				} else {
					log.Printf("[GetRole()] %s\n\n", errGetRole)
				}
			} else {
				output <- iamRole
			}
		}

		close(output)
	}()

	return
}

func genInstanceProfiles(getRole *iam.GetRoleOutput) (output chan *iam.ListInstanceProfilesForRoleOutput) {
	output = make(chan *iam.ListInstanceProfilesForRoleOutput)

	go func() {
		profilesForRole, errGetProfiles := svcIam.ListInstanceProfilesForRole(
			&iam.ListInstanceProfilesForRoleInput{
				RoleName: getRole.Role.RoleName,
			},
		)
		if errGetProfiles != nil {
			log.Println("[ListInstanceProfilesForRole()]", errGetProfiles)
		} else {
			output <- profilesForRole
		}

		close(output)
	}()

	return
}

type roleIPRemove struct {
	instanceProfilesOutput *iam.ListInstanceProfilesForRoleOutput
	roleOutput             *iam.GetRoleOutput
}

func removeInstanceProfilesFromRole(removeThese roleIPRemove) {
	for _, instanceProfile := range removeThese.instanceProfilesOutput.InstanceProfiles {
		removeRoleFromIP, errRemoveRoleFromIP := svcIam.RemoveRoleFromInstanceProfile(
			&iam.RemoveRoleFromInstanceProfileInput{
				InstanceProfileName: instanceProfile.InstanceProfileName,
				RoleName:            removeThese.roleOutput.Role.RoleName,
			},
		)
		if errRemoveRoleFromIP != nil {
			log.Println("[RemoveRoleFromInstanceProfile()]", errRemoveRoleFromIP)
		} else {
			fmt.Println(removeRoleFromIP)
		}
	}
}

func deleteRole(roleName string) {
	deleteRole, errDeleteRole := svcIam.DeleteRole(
		&iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		},
	)
	if errDeleteRole != nil {
		log.Printf("[DeleteRole()] role: '%s'\nerror: %s\n\n", roleName, errDeleteRole)
	} else {
		fmt.Printf("Deleted '%s': %s\n\n", roleName, deleteRole)
	}
}

func genRolePolicies(getRole *iam.GetRoleOutput) (output chan *iam.ListRolePoliciesOutput) {
	output = make(chan *iam.ListRolePoliciesOutput)

	go func() {
		policiesForRole, errGetRolePolicies := svcIam.ListRolePolicies(
			&iam.ListRolePoliciesInput{
				RoleName: getRole.Role.RoleName,
			},
		)
		if errGetRolePolicies != nil {
			log.Println("[ListRolePolicies()]", errGetRolePolicies)
		} else {
			output <- policiesForRole
		}

		close(output)
	}()

	return
}

type rolePolicyRemove struct {
	policiesOutput *iam.ListRolePoliciesOutput
	roleOutput     *iam.GetRoleOutput
}

func deletePoliciesFromRole(removeThese rolePolicyRemove) {
	for _, rolePolicy := range removeThese.policiesOutput.PolicyNames {
		deletePolicyFromRole, errDeletePolicyFromRole := svcIam.DeleteRolePolicy(
			&iam.DeleteRolePolicyInput{
				PolicyName: rolePolicy,
				RoleName:   removeThese.roleOutput.Role.RoleName,
			},
		)
		if errDeletePolicyFromRole != nil {
			log.Println("[DeleteRolePolicy()]", errDeletePolicyFromRole)
		} else {
			fmt.Println(deletePolicyFromRole)
		}
	}
}

func genAttachedPolicies(getRole *iam.GetRoleOutput) (output chan *iam.ListAttachedRolePoliciesOutput) {
	output = make(chan *iam.ListAttachedRolePoliciesOutput)

	go func() {
		policiesAttachedToRole, errPoliciesAttachedToRole := svcIam.ListAttachedRolePolicies(
			&iam.ListAttachedRolePoliciesInput{
				RoleName: getRole.Role.RoleName,
			},
		)
		if errPoliciesAttachedToRole != nil {
			log.Println("[ListAttachedRolePolicies()]", errPoliciesAttachedToRole)
		} else {
			output <- policiesAttachedToRole
		}

		close(output)
	}()

	return
}

type rolePolicyDetach struct {
	attachedPolicies *iam.ListAttachedRolePoliciesOutput
	roleOutput       *iam.GetRoleOutput
}

func detachPoliciesFromRole(removeThese rolePolicyDetach) {
	for _, rolePolicy := range removeThese.attachedPolicies.AttachedPolicies {
		detachPolicyFromRole, errDetachPolicyFromRole := svcIam.DetachRolePolicy(
			&iam.DetachRolePolicyInput{
				PolicyArn: rolePolicy.PolicyArn,
				RoleName:  removeThese.roleOutput.Role.RoleName,
			},
		)
		if errDetachPolicyFromRole != nil {
			log.Println("[DetachRolePolicy()]", errDetachPolicyFromRole)
		} else {
			fmt.Println(detachPolicyFromRole)
		}
	}
}
