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
import { showConfirm, showError } from '../../utils/alert';
import ToggleSwitch from '../ToggleSwitch/ToggleSwitch';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import btnStyles from '../../styles/button.module.scss';
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
  const [updated, setUpdated] = useState( false );

  // Modal Setup
  const [modalIsOpen, setModalIsOpen] = useState( false );

  // Modal Controls
  const openModal = () => setModalIsOpen( true );
  const closeModal = async ( skipConfirm = false ) => {
    // Clear the new team modal when closed.
    if ( !team ) {
      setLocalTeam( { ...localTeam, name: '', aprimoName: '' } );
      setUpdated( false );
    }

    let shouldClose = true;

    if ( updated && !skipConfirm ) {
      const { isConfirmed } = await showConfirm( 'Are you sure you want to close the edit window?  You will lose all unsaved progress.' );

      shouldClose = isConfirmed;
    }

    setModalIsOpen( !shouldClose );
  };

  // Update Controls
  const handleUpdate = ( key: keyof ITeam, value: ValueOf<ITeam> ) => {
    setLocalTeam( { ...localTeam, [key]: value } );
    setUpdated( true );
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
      const { data, error } = await response.json();

      newList = data;
      errMessage = error;
    } else {
      const response = await buildQuery( 'team', { active, team: id, teamName: name, teamAprimo: aprimoName }, 'PUT' );
      const { data, error } = await response.json();

      newList = data;
      errMessage = error;
    }

    // Update the team list with new data from the API.
    if ( newList ) {
      setTeams( newList );
    }

    if ( errMessage ) {
      showError( `Unable to complete your request. Reason: ${errMessage}` );
    } else {
      closeModal( true );
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
      <button className={ btnStyles['anchor-btn'] } onClick={ openModal } type="button">{ anchor }</button>
      <Modal
        isOpen={ modalIsOpen }
        onRequestClose={ () => closeModal() }
        contentLabel="Example Modal"
        style={ modalStyle }
      >
        <h2 className={ style.header }>
          { localTeam.id ? `Update ${localTeam.name}` : 'Add a New Team' }
        </h2>
        <label className={ style.label } htmlFor="team-name-input">
          Team Name
        </label>
        <input
          className={ style.input }
          id="team-name-input"
          type="text"
          value={ localTeam.name || '' }
          onChange={ e => handleUpdate( 'name', e.target.value ) }
          aria-label="Team Name"
        />
        <label className={ style.label } htmlFor="aprimo-name-input">
          Aprimo Name
        </label>
        <span className={ style.note }>
          This value must match the &quot;Team Name&quot; property
          <br />
          of an existing team in Content Commons.
        </span>
        <input
          className={ style.input }
          id="aprimo-name-input"
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
            className={ `${btnStyles.btn} ${btnStyles['spaced-btn']} ${updated ? '' : btnStyles['disabled-btn']}` }
            onClick={ handleSubmit }
            type="button"
            disabled={ !updated }
          >
            Submit
          </button>
          <button
            className={ `${btnStyles.btn} ${btnStyles['back-btn']} ` }
            onClick={ () => closeModal() }
            type="button"
          >
            Cancel
          </button>
        </div>
      </Modal>
    </>
  );
};
