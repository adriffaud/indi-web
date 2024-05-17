package fr.driffaud.indiweb;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;
import java.util.concurrent.CompletableFuture;

public class IndiClient {
    private Socket socket;
    private PrintWriter out;
    private BufferedReader in;

    public CompletableFuture<Void> start(String address, int port) throws IOException {
        socket = new Socket(address, port);
        out = new PrintWriter(socket.getOutputStream(), true);
        in = new BufferedReader(new InputStreamReader(socket.getInputStream()));

        return CompletableFuture.runAsync(this::listen);
    }

    private void listen() {
        String xmlElement;
        while (true) {
            try {
                if ((xmlElement = in.readLine()) == null) break;
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
            System.out.println("Received: " + xmlElement);
        }
    }

    public void getProperties() {
        this.sendMessage("<getProperties version=\"1.7\" />");
    }

    private void sendMessage(String msg) {
        System.out.println("Sending: " + msg);
        out.println(msg);
    }

    public void stopConnection() throws IOException {
        in.close();
        socket.close();
    }
}
