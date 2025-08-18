import './Input.scss';

const Input = ({ type, placeholder, value, onChange, name, readOnly }) => {
    return (
        <input
            className="custom-input"
            type={type}
            placeholder={placeholder}
            value={value}
            onChange={onChange}
            name={name}
            required
            readOnly={readOnly} // Bu satırı ekleyin
        />
    );
};

export default Input;