. $PSScriptRoot\Octopus.ps1

Describe 'Terraform Provider' {
    It 'must create a project group called Test' {
        Wait-ForOctopus
        $repository = Connect-ToOctopus http://localhost:8080
        $entity = Invoke-ScriptBlockWithRetries {$repository.ProjectGroups.FindByName("Test")}
        $entity | Should -Not -Be $null
    }
}
