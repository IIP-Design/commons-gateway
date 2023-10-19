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
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IBtnText {
  cancelButtonText?: string;
  confirmButtonText?: string;
  denyButtonText?: string;
}

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
export const ACCEPT_BTN_STYLE_CLASSES = `${btnStyles.btn} ${btnStyles['spaced-btn']}`;
export const NEUTRAL_BTN_STYLE_CLASSES = `${btnStyles['btn-light']} ${btnStyles['spaced-btn']}`;
export const REJECT_BTN_STYLE_CLASSES = `${btnStyles.btn} ${btnStyles['back-btn']} ${btnStyles['spaced-btn']}`;

// ////////////////////////////////////////////////////////////////////////////
// Exports
// ////////////////////////////////////////////////////////////////////////////
export const showError = ( text: string ) => Swal.fire( {
  icon: 'error',
  title: 'Input Error',
  text,
  customClass: {
    confirmButton: REJECT_BTN_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );

export const showSuccess = ( text: string ) => Swal.fire( {
  icon: 'success',
  title: 'Success!',
  text,
  customClass: {
    confirmButton: ACCEPT_BTN_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );

export const showInfo = ( title: string, text: string ) => Swal.fire( {
  icon: 'info',
  title,
  text,
  customClass: {
    confirmButton: ACCEPT_BTN_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );

export const showWarning = ( text: string, heading?: string ) => Swal.fire( {
  icon: 'warning',
  title: heading || 'Warning',
  text,
  customClass: {
    confirmButton: REJECT_BTN_STYLE_CLASSES,
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
    cancelButton: REJECT_BTN_STYLE_CLASSES,
    confirmButton: ACCEPT_BTN_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
} );

export const showTernary = ( text: string, buttons: IBtnText = {} ) => Swal.fire( {
  icon: 'info',
  title: 'Select an Option',
  text,
  showCancelButton: true,
  showConfirmButton: true,
  showDenyButton: true,
  customClass: {
    cancelButton: NEUTRAL_BTN_STYLE_CLASSES,
    confirmButton: ACCEPT_BTN_STYLE_CLASSES,
    denyButton: REJECT_BTN_STYLE_CLASSES,
    popup: alertStyles.text,
  },
  buttonsStyling: false,
  cancelButtonText: buttons.cancelButtonText || 'Cancel',
  confirmButtonText: buttons.confirmButtonText || 'Confirm',
  denyButtonText: buttons.denyButtonText || 'Deny',
} );
