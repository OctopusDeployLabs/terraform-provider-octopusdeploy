pushd tests

docker-compose up --detach

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

} until (Test-Connection -IPv4 -ComputerName localhost -TCPPort 8080 -Quiet)

popd
