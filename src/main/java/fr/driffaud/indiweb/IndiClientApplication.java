package fr.driffaud.indiweb;

import java.io.IOException;

public class IndiClientApplication {
    public static void main(String[] args) {
        try {
            var indiClient = new IndiClient();
            indiClient.start("localhost", 7624);

        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
