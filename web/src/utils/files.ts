// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import prettyBytes from 'pretty-bytes';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { submitFiles } from './upload';
import { showError, showSuccess } from './alert';

// ////////////////////////////////////////////////////////////////////////////
// Constants
// ////////////////////////////////////////////////////////////////////////////
const MAX_FILE_SIZE = 100 * 1000 * 1000;

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
let fileToUpload: File|null = null;

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const validateFile = ( { type, size }: File ) => {
  if ( !type.match( /^(image|video)\/.+/ ) ) {
    showError( 'Only pictures and/or videos may be uploaded' );

    return false;
  } if ( size > MAX_FILE_SIZE ) {
    showError( `Max file size is ${prettyBytes( MAX_FILE_SIZE )}` );

    return false;
  }

  return true;
};

const setUpload = ( file: File ) => {
  if ( !validateFile( file ) ) {
    return;
  }

  const list = document.getElementById( 'file-list' ) as HTMLElement;

  list.innerHTML = `${file.name} (${prettyBytes( file.size )})`;

  fileToUpload = file;
};

const handleFile = ( files?: FileList|null ) => {
  if ( files && files.length > 1 ) {
    showError( 'Only single-file uploads are currently supported' );
  } else if ( files ) {
    setUpload( files[0] );
  }
};

const validateSubmission = ( descriptionElem: HTMLInputElement ) => {
  const description = descriptionElem?.value;
  const file = fileToUpload;

  let error = false;

  // Error Check
  if ( !descriptionElem ) {
    showError( 'Internal error' );
    error = true;
  } else if ( !file ) {
    showError( 'No file has been selected for upload' );
    error = true;
  } else if ( !description ) {
    showError( 'No file description provided' );
    error = true;
  }

  return { description, file, error };
};

// ////////////////////////////////////////////////////////////////////////////
// Exports
// ////////////////////////////////////////////////////////////////////////////

/**
   * Prevents the default browser behavior (i.e. opening the file) when
   * a file is dropped into the browser.
   *
   * @param e The dragenter/dragover event.
   */
export const haltEvent = ( e: Event ) => {
  e.stopPropagation();
  e.preventDefault();
};

/**
 * Prepares the drag and dropped files for upload.
 *
 * @param e The drop event.
 */
export const dropHandler = ( e: DragEvent ) => {
  haltEvent( e );

  const files = e?.dataTransfer?.files;

  handleFile( files );
};

export const chooseHandler = ( e: Event ) => {
  const { files } = ( e.target as HTMLInputElement );

  handleFile( files );
};

export const submitHandler = async () => {
  // Prepare and validate
  const descriptionElem = document.getElementById( 'description-text' ) as HTMLInputElement;
  const fileElem = document.getElementById( 'file-list' ) as HTMLInputElement;

  const { file, description, error } = validateSubmission( descriptionElem );

  if ( error ) {
    return;
  }

  // Send data
  const response = await submitFiles( file as File, { description } );

  if ( response !== 'ok' ) {
    showError( 'Could not upload file' );
  } else {
    showSuccess( 'File has been uploaded' );
  }

  // Cleanup
  descriptionElem.value = '';
  fileElem.innerHTML = '';
  fileToUpload = null;
};
