import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Input from '../components/common/Input';
import Button from '../components/common/Button';
import { register } from '../api/authService';
import { FiMic } from "react-icons/fi"; // İkonu import ediyoruz
import { BsFillMicFill } from "react-icons/bs";
import './AuthPage.scss';

const RegisterPage = () => {
    const [formData, setFormData] = useState({ firstName: '', lastName: '', email: '', password: '' });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            await register(formData.firstName, formData.lastName, formData.email, formData.password);
            navigate('/verify-email', { state: { email: formData.email } });
        } catch (err) {
            setError(err.response?.data?.error || 'Kayıt sırasında bir hata oluştu.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-page">
            <div className="background-icons">
                <div className="background-icon-wrapper icon1"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon2"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon3"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon4"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon5"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon6"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon7"><BsFillMicFill /></div>
            </div>
            <div className="auth-container">
                <div className="auth-panel">
                    <h1>SpeeDh</h1>
                    <p>Fikirlerini sese, sesini metne dönüştür. Anında. Akıllıca.</p>
                </div>
                <div className="auth-form-wrapper">
                    <form className="auth-form" onSubmit={handleSubmit}>
                        <h2>Yeni Hesap Oluştur</h2>
                        <Input name="firstName" type="text" placeholder="Ad" value={formData.firstName} onChange={handleChange} />
                        <Input name="lastName" type="text" placeholder="Soyad" value={formData.lastName} onChange={handleChange} />
                        <Input name="email" type="email" placeholder="E-posta Adresi" value={formData.email} onChange={handleChange} />
                        <Input name="password" type="password" placeholder="Şifre" value={formData.password} onChange={handleChange} />
                        <Button type="submit" disabled={loading}>
                            {loading ? 'Oluşturuluyor...' : 'Kayıt Ol'}
                        </Button>
                        {error && <p className="error-message">{error}</p>}
                        <p className="form-footer">
                            Zaten bir hesabın var mı? <a href="/login">Giriş Yap</a>
                        </p>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default RegisterPage;