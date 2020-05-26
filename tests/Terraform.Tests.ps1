. $PSScriptRoot\Octopus.ps1

Describe 'Terraform Provider' {
    It 'project-group.tf must create a project group called Test' {
        $repository = Connect-ToOctopus http://localhost:8080
        $entity = Invoke-ScriptBlockWithRetries {$repository.ProjectGroups.FindByName("Test")}
        $entity | Should -Not -Be $null
        $entity.Name | Should -Be "Test"
        $entity.Description | Should -Be "Test Applications"
    }
}
