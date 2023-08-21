// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import prettyBytes from 'pretty-bytes';
import Swal from 'sweetalert2';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { submitFiles } from './api';

// ////////////////////////////////////////////////////////////////////////////
// Constants
// ////////////////////////////////////////////////////////////////////////////
const MAX_FILE_SIZE = 100 * 1000 * 1000;

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const formData = new FormData();

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const showError = ( text: string ) => Swal.fire( {
  icon: 'error',
  title: 'Input Error',
  text,
} );

const showSuccess = ( text: string ) => Swal.fire( {
  icon: 'success',
  title: 'Success!',
  text,
} );

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

const addToUploadList = ( file: File ) => {
  if ( !validateFile( file ) ) {
    return;
  }

  const list = document.getElementById( 'file-list' );
  const listItem = document.createElement( 'li' );

  listItem.innerHTML = `${file.name} (${prettyBytes( file.size )})`;
  list?.appendChild( listItem );

  formData.append( 'file', file );
};

const handleFiles = ( files: FileList ) => {
  [...files].forEach( file => addToUploadList( file ) );
};

const validateSubmission = ( descriptionElem: HTMLInputElement, listElem: HTMLElement|null ) => {
  const description = descriptionElem?.value;

  const listEntries = listElem?.childElementCount;
  const plural = ( listEntries && listEntries > 1 ) ? 's' : '';

  let totalSizeBytes = 0;

  formData.forEach( entry => {
    console.log( entry );
    totalSizeBytes += ( entry as File ).size;
  } );

  const ret = { description, listEntries, totalSizeBytes, plural, error: false };

  // Error Check
  if ( !descriptionElem || !listElem ) {
    showError( 'Internal error' );
    ret.error = true;
  } else if ( !listEntries ) {
    showError( 'No files have been selected for upload' );
    ret.error = true;
  } else if ( !description ) {
    showError( 'No file description provided' );
    ret.error = true;
  }

  if ( totalSizeBytes > MAX_FILE_SIZE ) {
    showError( `File${plural} total size is ${prettyBytes( totalSizeBytes )}, but the maximum allowed is ${prettyBytes( MAX_FILE_SIZE )}` );
    ret.error = true;
  }

  return ret;
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

  if ( files ) {
    handleFiles( files );
  }
};

export const chooseHandler = ( e: Event ) => {
  const { files } = ( e.target as HTMLInputElement );

  if ( files ) {
    handleFiles( files );
  }
};

export const submitHandler = async () => {
  // Prepare and validate
  const descriptionElem = document.getElementById( 'description-text' ) as HTMLInputElement;
  const listElem = document.getElementById( 'file-list' ) as HTMLElement;
  const { description, totalSizeBytes, plural, error } = validateSubmission( descriptionElem, listElem );

  if ( error ) {
    return;
  }
  formData.set( 'description', description );

  // Send data
  const response = await submitFiles( 'upload', formData );

  if ( !response.ok ) {
    showError( `Error on the ${response.status >= 500 ? 'server' : 'client'}` );
  } else {
    showSuccess( `File${plural} have been uploaded` );
  }

  // Cleanup
  descriptionElem.value = '';
  listElem.innerHTML = '';
  formData.delete( 'file' );
  formData.delete( 'description' );
};
