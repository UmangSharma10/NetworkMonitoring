import net.minidev.json.parser.ParseException;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.*;


public class Main {
    public static void functionProcess(String Credentials){
        try {
            List<String> commands = new ArrayList<>();
            commands.add("/home/umang/GolandProjects/sample/plugin.exe");
           // String encodedString = Base64.getEncoder().encodeToString(jsonArray.toString().getBytes(StandardCharsets.UTF_8));


            commands.add(Credentials);
            ProcessBuilder processBuilder = new ProcessBuilder(commands);
            Process process = processBuilder.start();

            // for reading the output from stream
            BufferedReader stdInput = new BufferedReader(new InputStreamReader(
                    process.getInputStream()));
            String s;
            while ((s = stdInput.readLine()) != null) {
                System.out.println(s);

            }
        } catch (IOException e) {
            e.printStackTrace();
        }

    }


    public static void main(String[] args) throws IOException, ParseException {
       new Bootstrap();
    }
}