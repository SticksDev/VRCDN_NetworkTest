Write-Host "Downloading and running VRCDN_NetworkTest.exe..."

# Define the GitHub repository and release information
$repoOwner = "SticksDev"
$repoName = "VRCDN_NetworkTest"
$releaseUrl = "https://api.github.com/repos/$repoOwner/$repoName/releases/latest"

# Get the latest release information
$releaseInfo = Invoke-RestMethod -Uri $releaseUrl

# Get the current architecture
$architecture = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "386" }

Write-Host "Got latest release information. Downloading..."

# Construct download URLs
$releaseFileName = "VRCDN_NetworkTest-1.0.0-windows-$architecture.zip"
$releaseFileUrl = $releaseInfo.assets | Where-Object { $_.name -eq $releaseFileName } | Select-Object -ExpandProperty browser_download_url
$md5FileUrl = "$releaseFileUrl.md5"

# Download release files
Invoke-WebRequest -Uri $releaseFileUrl -OutFile $releaseFileName
Invoke-WebRequest -Uri $md5FileUrl -OutFile "$releaseFileName.md5"

Write-Host "Downloaded release files. Verifying hashes..."

# Verify hashes
$downloadedHash = Get-FileHash -Algorithm MD5 -Path $releaseFileName | Select-Object -ExpandProperty Hash
$expectedHash = Get-Content "$releaseFileName.md5"

if ($downloadedHash -eq $expectedHash) {
    Write-Host "Hash verification successful. Unzipping..."
}
else {
    Write-Host "Hash verification failed. Exiting script."
    Remove-Item -Path $releaseFileName, "$releaseFileName.md5" -Force
    Exit
}

# Unzip the downloaded file
Expand-Archive -Path $releaseFileName -DestinationPath .

Write-Host "Unzipped. We will now run VRCDN_NetworkTest.exe. When you are done, press any key to continue after the program exits."

# Run VRCDN_NetworkTest.exe (in our current powershell session)
.\VRCDN_NetworkTest.exe

# Wait for keypress
Write-Host "Press any key to continue..."
$null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')

# Clean up: Delete the zip, md5, and exe files
Remove-Item -Path $releaseFileName, "$releaseFileName.md5", ".\VRCDN_NetworkTest.exe" -Force
