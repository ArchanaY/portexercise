$env:BASE_URL="http://127.0.0.1:8080"

# Heartbeat
Write-Output "Succeed heartbeat 200"
Invoke-WebRequest http://127.0.0.1:8080/heartbeat

Write-Output "succeed POST returns 201"
Invoke-WebRequest $env:BASE_URL/port -Method Post -InFile .\tests\example1.json -Headers @{'Content-Type' = 'application/json; charset=utf-8'}

Write-Output "Succeed GET 200"
Invoke-WebRequest $env:BASE_URL/port/key2 -Method Get -OutFile .\tests\key2.json

Write-Output "succeed POST returns 201"
Invoke-WebRequest $env:BASE_URL/port -Method Post -InFile .\tests\example2.json -Headers @{'Content-Type' = 'application/json; charset=utf-8'}

Write-Output "Succeed GET 200"
Invoke-WebRequest $env:BASE_URL/port/key2 -Method Get -OutFile .\tests\updated.json

Write-Output "succeed POST returns 201"
Invoke-WebRequest $env:BASE_URL/port -Method Post -InFile .\tests\ports.json -Headers @{'Content-Type' = 'application/json; charset=utf-8'}

Write-Output "Succeed GET 200"
Invoke-WebRequest $env:BASE_URL/port/ZWHRE -Method Get -OutFile .\tests\port_out.json
