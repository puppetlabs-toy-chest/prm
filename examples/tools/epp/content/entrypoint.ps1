$Files = If ([string]::IsNullOrEmpty($args)) {
  Get-ChildItem -Recurse -File -Filter '*.epp'
} ElseIf ($args[0].ToString() -eq 'help') {
  "Please specify a filter pattern to search for EPP files, such as '*.epp'"
  return
} Else {
  $args | ForEach-Object {
    Get-ChildItem -Recurse -File -Filter $_
  }
}
$Files = $Files.FullName

"Validating EPP for ${files}"
& puppet epp validate --continue_on_error $Files
