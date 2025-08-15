import './Button.scss';

const Button = ({ children, type = 'button', onClick, disabled }) => {
    return (
        <button
            className="custom-button"
            type={type}
            onClick={onClick}
            disabled={disabled}
        >
            {children}
        </button>
    );
};

export default Button;