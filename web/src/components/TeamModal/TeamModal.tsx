// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useState } from 'react';
import type{ FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import Modal from 'react-modal';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { buildQuery } from '../../utils/api';
import { showError } from '../../utils/alert';
import ToggleSwitch from '../ToggleSwitch/ToggleSwitch';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import btnStyle from '../../styles/button.module.scss';
import style from './TeamModal.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface ITeamModalProps {
  readonly team?: ITeam;
  readonly setTeams: React.Dispatch<React.SetStateAction<ITeam[]>>;
  readonly anchor: string | JSX.Element;
}

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
Modal.setAppElement( document.getElementById( 'root' ) as HTMLElement );

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const TeamModal: FC<ITeamModalProps> = ( { team, setTeams, anchor }: ITeamModalProps ) => {
  // Team Setup
  const [localTeam, setLocalTeam] = useState<Partial<ITeam>>( team || { active: true } );

  // Modal Setup
  const [modalIsOpen, setModalIsOpen] = useState( false );

  // Modal Controls
  const openModal = () => setModalIsOpen( true );
  const closeModal = () => {
    // Clear the new team modal when closed.
    if ( !team ) {
      setLocalTeam( { ...localTeam, name: '', aprimoName: '' } );
    }
    setModalIsOpen( false );
  };

  // Update Controls
  const handleUpdate = ( key: keyof ITeam, value: ValueOf<ITeam> ) => {
    setLocalTeam( { ...localTeam, [key]: value } );
  };

  // Update Team
  const handleSubmit = async () => {
    const { name, aprimoName, id, active } = localTeam;

    if ( !name ) {
      showError( 'A team must have a name' );

      return;
    }

    if ( !aprimoName ) {
      showError( 'Please specify the value Aprimo uses for this team' );

      return;
    }

    let newList;
    let errMessage;

    // If the team is new, send a create request, otherwise send an update request.
    if ( !id ) {
      const response = await buildQuery( 'team', { teamName: name, teamAprimo: aprimoName }, 'POST' );
      const { data, message } = await response.json();

      newList = data;
      errMessage = message;
    } else {
      const response = await buildQuery( 'team', { active, team: id, teamName: name, teamAprimo: aprimoName }, 'PUT' );
      const { data, message } = await response.json();

      newList = data;
      errMessage = message;
    }

    // Update the team list with new data from the API.
    if ( newList ) {
      setTeams( newList );
    }

    if ( errMessage ) {
      showError( `Unable to complete your request. Reason: ${errMessage}` );
    } else {
      closeModal();
    }
  };

  // Styles
  const modalStyle = {
    content: {
      height: 'fit-content',
      width: '400px',
      minWidth: 'fit-content',
      top: '50%',
      left: '50%',
      right: 'auto',
      bottom: 'auto',
      marginRight: '-50%',
      padding: '2rem',
      transform: 'translate(-50%, -50%)',
      boxShadow: 'rgba(0, 0, 0, 0.35) 0px 5px 15px',
    },
  };

  return (
    <>
      <button className={ btnStyle['anchor-btn'] } onClick={ openModal } type="button">{ anchor }</button>
      <Modal
        isOpen={ modalIsOpen }
        onRequestClose={ closeModal }
        contentLabel="Example Modal"
        style={ modalStyle }
      >
        <h2 className={ style.header }>
          { localTeam.id ? `Update ${localTeam.name}` : 'Add a New Team' }
        </h2>
        <label className={ style.label }>
          Team Name
        </label>
        <input
          className={ style.input }
          type="text"
          value={ localTeam.name || '' }
          onChange={ e => handleUpdate( 'name', e.target.value ) }
          aria-label="Team Name"
        />
        <label className={ style.label }>
          Aprimo Name
        </label>
        <input
          className={ style.input }
          type="text"
          value={ localTeam.aprimoName || '' }
          onChange={ e => handleUpdate( 'aprimoName', e.target.value ) }
          aria-label="Team Name"
        />
        { localTeam.id && (
          <div className={ style.label }>
            <ToggleSwitch
              active={ localTeam.active ?? false }
              callback={ e => handleUpdate( 'active', e ) }
              id={ localTeam.id }
            />
          </div>
        ) }
        <div className={ style['btn-container'] }>
          <button
            className={ btnStyle.btn }
            onClick={ handleSubmit }
            type="button"
          >
            Submit
          </button>
          <button
            className={ `${btnStyle.btn} ${btnStyle['back-btn']} ` }
            onClick={ closeModal }
            type="button"
          >
            Cancel
          </button>
        </div>
      </Modal>
    </>
  );
};
