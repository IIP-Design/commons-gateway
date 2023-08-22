// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import Swal from 'sweetalert2';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import btnStyles from '../styles/button.module.scss';
import alertStyles from '../styles/alert.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const ACCEPT_STYLE_CLASSES = `${btnStyles.btn} ${btnStyles['spaced-btn']}`;
const REJECT_STYLE_CLASSES = `${btnStyles.btn} ${btnStyles['back-btn']} ${btnStyles['spaced-btn']}`;

// ////////////////////////////////////////////////////////////////////////////
// Exports
// ////////////////////////////////////////////////////////////////////////////
export const showError = ( text: string ) => Swal.fire( {
  icon: 'error',
  title: 'Input Error',
  text,
  customClass: {
    confirmButton: REJECT_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );

export const showSuccess = ( text: string ) => Swal.fire( {
  icon: 'success',
  title: 'Success!',
  text,
  customClass: {
    confirmButton: ACCEPT_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );

export const showConfirm = ( text: string ) => Swal.fire( {
  icon: 'warning',
  title: 'Please confirm your selection',
  text,
  showCancelButton: true,
  showConfirmButton: true,
  customClass: {
    confirmButton: ACCEPT_STYLE_CLASSES,
    cancelButton: REJECT_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );
