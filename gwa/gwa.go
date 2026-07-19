package main

/*
# usage: git worktree add [-f] [--detach] [--checkout] [--lock [--reason <string>]]
#                         [--orphan] [(-b | -B) <new-branch>] <path> [<commit-ish>]
#
#     -f, --[no-]force      checkout <branch> even if already checked out in other worktree
#     -b <branch>           create a new branch
#     -B <branch>           create or reset a branch
#     --[no-]orphan         create unborn branch
#     -d, --[no-]detach     detach HEAD at named commit
#     --[no-]checkout       populate the new working tree
#     --[no-]lock           keep the new working tree locked
#     --[no-]reason <string>
#                           reason for locking
#     -q, --[no-]quiet      suppress progress reporting
#     --[no-]track          set up tracking mode (see git-branch(1))
#     --[no-]guess-remote   try to match the new branch name with a remote-tracking branch

# git worktree add --help
# git worktree add -b DEVPL-442/bau/crossplane-version-bumps ../DEVPL-442_bau_crossplane-version-bumps
# git worktree add -b DEVPL-567/crossplane/gcp-provider/upgrade/v1.8.0 ../xp-gcp-v1.8.0
# git worktree add -b DEVPL-567/crossplane/upgrade/v1.17.0 ../xp-v1.17.0
# git worktree add -b DEVPL-567/gcp-provider/patch-1.8.1 ../DEVPL-567_gcp-provider_patch-1.8.1
# git worktree add -b DEVPL-567/hack-crossplane/fix-var-file-paths ../DEVPL-567_hack-crossplane_fix-var-file-paths
# git worktree add -b DEVPL-567/hack-crossplane/unique-names ../DEVPL-567_hack-crossplane_unique-names origin/HEAD
# git worktree add -b DEVPL-655/crossplane-into-helm-charts ../devpl-655
# git worktree add -b DEVPL-656/hack-xp-taskfile/xrd-comp-picker ../DEVPL-656_hack-xp-taskfile_xrd-comp-picker
# git worktree add -b DEVPL-656/new-resource-howto ../DEVPL-656_new-resource-howto
# git worktree add -b DEVPL-667/refactor-composition-mode ../../../devpl-667
# git worktree add -b DEVPL-684/fit-and-finish/odds-and-sods ../devpl-684
# git worktree add -b DEVPL-684/helm-xp-r/dry ../DEVPL-684_helm-xp-r_dry
# git worktree add -b DEVPL-708/helm-xp-providers/gcp-db/dev ../DEVPL-708_helm-xp-providers_gcp-db_dev
# git worktree add -b DEVPL-708/helm-xp-resources/set-gcp-prov-cfg ../DEVPL-708_helm-xp-resources_set-gcp-prov-cfg
# git worktree add -b DEVPL-708/xp-db/dev ../DEVPL-708_xp-db_dev
# git worktree add -b NOJIRA/hack-xp/doctor-etc ../NOJIRA_hack-xp_doctor-etc
# git worktree add -b NOJIRA/hack-xp/doctor-etc2 ../NOJIRA_hack-xp_doctor-etc2
# git worktree add -b NOJIRA/helm-xp-r/rearrange-dirs ../NOJIRA_helm-xp-r_rearrange-dirs
# git worktree add -b rebuild/simple-debian-vm ../rebuild_simple-debian-vm origin/HEAD
# git worktree add ../_factorio-server-kit origin/main
# git worktree add ../DEVPL-442_bau_crossplane-version-bumps DEVPL-442/bau/crossplane-version-bumps
# git worktree add ../DEVPL-567_hack-crossplane_fix-var-file-paths DEVPL-567/hack-crossplane/fix-var-file-paths
# git worktree add ../DEVPL-656_hack-xp-taskfile_xrd-comp-picker/
# git worktree add ../DEVPL-670_crossplane-storage-test DEVPL-670/crossplane-storage-test
# git worktree add ../DEVPL-670_crossplane-storage-test origin/DEVPL-670/crossplane-storage-test
# git worktree add ../DEVPL-698_dynamodb-crossplane-implementation DEVPL-698/dynamodb-crossplane-implementation
# git worktree add ../NOJIRA_hack-xp_doctor-etc2 NOJIRA/hack-xp/doctor-etc2
# git worktree add ../pr-1145 DEVPL-630/ArgoIngres-Review
# git worktree add ../xp-v1.17.0 DEVPL-567/crossplane/upgrade/v1.17.0
# git worktree add DEVPL-567/crossplane/upgrade/v1.17.0 ../xp-v1.17.0

# git worktree remove --force ../NOJIRA_hack-xp_doctor-etc/
# git worktree remove .
# git worktree remove ../../../devpl-655/
# git worktree remove ../../../DEVPL-708_helm-xp-resources_set-gcp-prov-cfg/
# git worktree remove ../DEVPL-442_bau_crossplane-version-bumps/
# git worktree remove ../DEVPL-567_gcp-provider_patch-1.8.1/
# git worktree remove ../DEVPL-567_hack-crossplane_fix-var-file-paths/
# git worktree remove ../DEVPL-567_hack-crossplane_unique-names/
# git worktree remove ../devpl-655/
# git worktree remove ../DEVPL-656_new-resource-howto/
# git worktree remove ../devpl-667/
# git worktree remove ../DEVPL-670_crossplane-storage-test/
# git worktree remove ../DEVPL-684_helm-xp-r_dry/
# git worktree remove ../DEVPL-708_helm-xp-providers_gcp-db_dev/
# git worktree remove ../DEVPL-708_xp-db_dev/
# git worktree remove ../NOJIRA_hack-xp_doctor-etc/
# git worktree remove ../NOJIRA_hack-xp_doctor-etc2/
# git worktree remove ../NOJIRA_helm-xp-r_rearrange-dirs/
# git worktree remove ../pr-1145/
# git worktree remove ../xp-v1.17.0/
*/
