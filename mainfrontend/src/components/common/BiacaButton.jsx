import './BiacaButton.scss';

const BiacaButton = () => {
    return (
        <a
            href="https://biacasoftware.com"
            target="_blank"
            rel="noopener noreferrer"
            className="biaca-button"
        >
            <span>Created By</span>
            &copy; <strong>Biaca Software</strong>
        </a>
    );
};

export default BiacaButton;