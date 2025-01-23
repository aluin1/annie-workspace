param (
    [string]$folderSource,
    [string]$folderSourceOut,
    [string]$folderTargetTxt,
    [string]$folderTargetGrp,
    [string]$folderDebug
)

# Menentukan platform (Windows atau Linux)
if ($IsWindows) {
    # Jalur untuk Windows
    $folderSource = "C:\Script\01-DS_RUN_Log"
    $folderSourceOut = "C:\Script\01-DS_RUN_Log\log"
    $folderTargetTxt = "C:\Script\02-Extract-Text-DB\input_text"
    $folderTargetGrp = "C:\Script\03-Extract-Graph-DB\1-input"
    $folderDebug = "C:\Script\02-Extract-Text-DB\output"
} else {
    # Jalur untuk Linux
    $folderSource = "/var/app_ds/01-DS_RUN_Log"
    $folderSourceOut = "/var/app_ds/01-DS_RUN_Log/log"
    $folderTargetTxt = "/var/app_ds/02-Extract-Text-DB/input_text"
    $folderTargetGrp = "/var/app_ds/03-Extract-Graph-DB/1-input"
    $folderDebug = "/var/app_ds/02-Extract-Text-DB/output"
}

# Fungsi untuk memeriksa apakah skrip sudah berjalan
function Test-IfAlreadyRunning {
    [CmdletBinding()]
    Param (
        [Parameter(Mandatory = $true)]
        [ValidateNotNullOrEmpty()]
        [String]$ScriptName
    )

    # Gunakan 'ps' untuk mendapatkan daftar proses pada Linux
    $processes = & ps -eo pid,cmd --no-headers | Out-String

    # Cek apakah ada skrip yang sama sedang berjalan
    $processes -split "`n" | ForEach-Object {
        $line = $_ -split '\s+'
        $otherPid = $line[0]  # Mengganti nama variabel untuk menghindari konflik dengan $PID
        $cmdline = $line[1..($line.Length - 1)] -join ' '

        If (($cmdline -match $ScriptName) -And ($otherPid -ne $PID)) {
            Write-Host "PID [$otherPid] is already running this script [$ScriptName]"
            Write-Host "Exiting this instance. (PID=[$PID])..."
            Exit
        }
    }
}

# Nama skrip saat ini
$ScriptName = $MyInvocation.MyCommand.Name
Test-IfAlreadyRunning -ScriptName $ScriptName

# Proses utama
Get-ChildItem -Path $folderSource -Filter *.case.txt | ForEach-Object {
    Write-Host ("Processing " + $_.FullName)

    # Buat target file
    $target1 = Join-Path $folderTargetTxt "Extract_iText_" + $_.Name.Replace(".case", "")
    $target2 = Join-Path $folderTargetGrp "Extract_Graph_" + $_.Name.Replace(".case", "")

    # Salin file
    Copy-Item $_.FullName $target1
    Copy-Item $_.FullName $target2

    # Filter konten file dan simpan kembali
    Set-Content -Path $target1 -Value (Get-Content -Path $target1 | Select-String -Pattern 'PLOTDRAW' -NotMatch)
    Set-Content -Path $target2 -Value (Get-Content -Path $target2 | Select-String -Pattern 'PLOTDRAW','1P0','1L0')

    # Tampilkan informasi file
    Write-Host ('File 1 : ' + $target1)
    Write-Host ('File 2 : ' + $target2)

    # Pindahkan file asli ke folder log
    Move-Item $_.FullName -Force -Destination $folderSourceOut
}

Write-Host "Proses selesai."
