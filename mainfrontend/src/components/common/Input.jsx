import './Input.scss';

const Input = ({ type, placeholder, value, onChange, name }) => {
    return (
        <input
            className="custom-input"
            type={type}
            placeholder={placeholder}
            value={value}
            onChange={onChange}
            name={name}
            required
        />
    );
};

export default Input;