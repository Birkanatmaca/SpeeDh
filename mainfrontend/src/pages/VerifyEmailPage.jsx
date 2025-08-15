import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom'; // Yönlendirme ve state için hook'lar
import Input from '../components/common/Input';
import Button from '../components/common/Button';
import { verifyEmail } from '../api/authService'; // Yeni API fonksiyonumuz
import { BsFillMicFill } from "react-icons/bs";
import './AuthPage.scss';

const VerifyEmailPage = () => {
    const [code, setCode] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const navigate = useNavigate(); // Başarılı olunca yönlendirme yapmak için
    const location = useLocation(); // Register sayfasından gönderilen e-posta bilgisini almak için


    const email = location.state?.email;


    if (!email) {
        return (
            <div className="auth-page">
                <p>Doğrulama için bir e-posta adresi bulunamadı. Lütfen önce kayıt olun.</p>
                <a href="/register">Kayıt Sayfasına Dön</a>
            </div>
        );
    }

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            const response = await verifyEmail(email, code);
            console.log('Hesap doğrulandı! Token:', response.data.token);
            alert('Hesabınız başarıyla doğrulandı! Giriş yaptınız.');

            // Gerçek bir uygulamada burada token'ı kaydedip ana sayfaya yönlendiririz:
            // navigate('/dashboard');

        } catch (err) {
            setError(err.response?.data?.error || 'Doğrulama sırasında bir hata oluştu.');
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
                    <h1>Son Bir Adım...</h1>
                    <p>E-posta adresinize gönderdiğimiz 6 haneli kodu girerek hesabınızı aktive edin.</p>
                </div>
                <div className="auth-form-wrapper">
                    <form className="auth-form" onSubmit={handleSubmit}>
                        <h2>Hesabı Doğrula</h2>
                        <p className="form-footer" style={{ marginBottom: '2rem' }}>
                            Doğrulama kodu <strong>{email}</strong> adresine gönderildi.
                        </p>
                        <Input
                            name="code" type="text" placeholder="6 Haneli Kod"
                            value={code} onChange={(e) => setCode(e.target.value)}
                        />
                        <Button type="submit" disabled={loading}>
                            {loading ? 'Doğrulanıyor...' : 'Doğrula ve Giriş Yap'}
                        </Button>
                        {error && <p className="error-message">{error}</p>}
                    </form>
                </div>
            </div>
        </div>
    );
};

export default VerifyEmailPage;