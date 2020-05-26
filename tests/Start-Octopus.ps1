pushd tests

docker-compose up --detach --quiet-pull --always-recreate-deps --force-recreate

$start = Get-Date
do {
    Write-Host "Waiting for Octopus"
    sleep 5
    $now = Get-Date
    $wait = New-Timespan -Start $start -End $now
    if ($wait.TotalMinutes -ge 5) {
        Write-Host "Gave up waiting"
        break;
    }

} until (Test-NetConnection -ComputerName localhost -Port 8080 -InformationLevel Quiet)

popd
