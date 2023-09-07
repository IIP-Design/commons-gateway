////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import type { FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { showConfirm } from '../utils/alert';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import styles from '../styles/button.module.scss'

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IBackButtonProps {
    id?: string;
    text?: string;
    showConfirmDialog?: boolean;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
const BackButton: FC<IBackButtonProps> = ( { id, text, showConfirmDialog }: IBackButtonProps ) => {
    const goBack = () => {
        if( showConfirmDialog ) {
            showConfirm( "Are you sure you want to return to the previous page?  You will lose all unsaved progress." )
                .then( ( result ) => {
                    if( result.isConfirmed ) {
                        window.history.back();
                    }
                } );
        }
    }

    return <button id={ id || "back-btn"} type="button" onClick={goBack} className={`${styles.btn} ${styles['back-btn']} ${styles['spaced-btn']}`}>{ text || "Back"}</button>
}

export default BackButton;