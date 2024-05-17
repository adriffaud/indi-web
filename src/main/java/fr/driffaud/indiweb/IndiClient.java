package fr.driffaud.indiweb;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;

public class IndiClient {
    private Socket socket;
    private PrintWriter out;
    private BufferedReader in;

    public void start(String address, int port) throws IOException {
        socket = new Socket(address, port);
        out = new PrintWriter(socket.getOutputStream(), true);
        in = new BufferedReader(new InputStreamReader(socket.getInputStream()));

        String xmlElement;
        while ((xmlElement = in.readLine()) != null) {
            System.out.println("Received: " + xmlElement);
        }
    }

    public void sendMessage(String msg) {
        out.println(msg);
    }

    public void stopConnection() throws IOException {
        in.close();
        socket.close();
    }
}
