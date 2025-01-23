#!/bin/bash
# Clear the terminal
clear
# Print report info
echo "Printing Report $1"
echo "$(date)"
echo "----------------------------------------"
# Run the RUNCASE program
${DIR_APP}/01-DS_RUN_Log/RUNCASE.exe $1 .
sleep 5
# Run DS_RUN program
${DIR_APP}/01-DS_RUN_Log/DS_RUN.exe Log $1
sleep 2
# Separate the report
echo "Separating Report $1"
echo "$(date)"
echo "----------------------------------------"
# chmod +r ${DIR_APP}/01-DS_RUN_Log/CSP-Extract-Data.ps1
# Run PowerShell script (adjust if needed, assuming pwsh is installed)
pwsh -ExecutionPolicy Unrestricted -File ${DIR_APP}/01-DS_RUN_Log/CSP-Extract-Data.ps1
# Move the .common.txt files
mv -f ${DIR_APP}/01-DS_RUN_Log/*.common.txt ${DIR_APP}/03-Extract-Graph-DB/1-input
sleep 2
# Run GetDataFromText program
${DIR_APP}/02-Extract-Text-DB/GetDataFromText.exe
sleep 2
# Completing the report
echo "Completing Report $1"
echo "$(date)"
echo "----------------------------------------"
# Move input text and graphic files to output directory
mv -f ${DIR_APP}/02-Extract-Text-DB/input_text/* ${DIR_APP}/02-Extract-Text-DB/output
mv -f ${DIR_APP}/02-Extract-Text-DB/input_graphic/* ${DIR_APP}/02-Extract-Text-DB/output
# Delete the Flag file
echo "Deleting Flag File"
rm -f ${DIR_APP}/01-DS_RUN_Log/AutoPrint.FLG
