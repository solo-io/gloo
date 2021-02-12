export function doDownload(downloadContent: string, filename: string) {
  const templElement = document.createElement("a");
  const file = new Blob([downloadContent], {
    type: "text/plain",
  });
  templElement.href = URL.createObjectURL(file);
  templElement.download = filename;
  document.body.appendChild(templElement); // Required for this to work in FireFox
  templElement.click();
}
