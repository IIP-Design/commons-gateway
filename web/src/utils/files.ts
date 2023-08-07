/**
   * Prevents the default browser behavior (i.e. opening the file) when
   * a file is dropped into the browser.
   *
   * @param e The dragenter/dragover event.
   */
export const dragHandler = (e: DragEvent) => {
  e.stopPropagation();
  e.preventDefault();
};

const addToUploadList = (file: string) => {
  const list = document.getElementById("file-list");
  const listItem = document.createElement("li");
  
  listItem.innerHTML = file;
  list?.appendChild(listItem);
}

const handleFiles = (files: FileList) => {
  [...files].forEach(file => addToUploadList(file.name));
};

/**
 * Prepares the drag and dropped files for upload.
 *
 * @param e The drop event.
 */
export const dropHandler = (e: DragEvent) => {
  e.stopPropagation();
  e.preventDefault();

  const files = e?.dataTransfer?.files;

  if (files) {
    handleFiles(files);
  }
};