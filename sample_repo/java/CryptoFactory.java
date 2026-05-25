import java.security.KeyPairGenerator;

public class CryptoFactory {
    public KeyPairGenerator generator() throws Exception {
        return KeyPairGenerator.getInstance("RSA");
    }
}
