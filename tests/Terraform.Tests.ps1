. $PSScriptRoot\Octopus.ps1

$repository = Connect-ToOctopus http://localhost:8080

Describe 'Terraform Provider' {
    It 'project-group.tf must create a project group called Test' {
        $entity = Invoke-ScriptBlockWithRetries {$repository.ProjectGroups.FindByName("Test")}
        $entity | Should -Not -Be $null
        $entity.Name | Should -Be "Test"
        $entity.Description | Should -Be "Test Applications"
    }

    It 'environment.tf must create an environment called TestEnv1' {
        $entity = Invoke-ScriptBlockWithRetries {$repository.Environments.FindByName("TestEnv1")}
        $entity | Should -Not -Be $null
        $entity.Name | Should -Be "TestEnv1"
    }
}
