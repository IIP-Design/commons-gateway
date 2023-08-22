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
import '../styles/button.scss'

interface IBackButtonProps {
    id?: string;
    showConfirmDialog?: boolean;
}

const BackButton: FC<IBackButtonProps> = ( { id, showConfirmDialog }: IBackButtonProps ) => {
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

    return <button id={ id || "back-btn"} type="button" onClick={goBack} className="back-btn">Back</button>
}

export default BackButton;