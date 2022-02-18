
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import net.minidev.json.parser.JSONParser;
import net.minidev.json.parser.ParseException;
import org.json.*;
import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.*;

import static scala.util.parsing.json.JSON.jsonArray;


public class Bootstrap{

        Bootstrap() throws IOException, ParseException {
            //Timer Function
            HashMap<String,String> ipAndTime= new HashMap<>();
            HashMap<String,String> OrignalTimer = new HashMap<String,String>();
            HashMap<String,String> Context = new HashMap<String,String>();
            List<Map<String,String>> contextMapIntoList = new ArrayList<>();


            JSONParser jsonParser = new JSONParser();
            try (FileReader reader = new FileReader("/home/umang/IdeaProjects/Plugins/src/main/java/cred.json"))
            {
                Object obj = jsonParser.parse(reader);
                JSONArray jsonArray = new JSONArray(obj.toString());


                for(int i = 0; i<jsonArray.length(); i++) {
                    String temp = jsonArray.get(i).toString();
                    ObjectMapper mapper = new ObjectMapper();
                    HashMap<String,String> tempMap = mapper.readValue(temp, new TypeReference<HashMap<String, String>>() {});
                    contextMapIntoList.add(tempMap);
                }

                int k = 0;
                for(Map<String, String> temp:contextMapIntoList)
                {
                    String ip = temp.get("host");
                    String time = temp.get("scheduleTime");
                    String encodedJsonStringARG1 = Base64.getEncoder().encodeToString(jsonArray.get(k).toString().getBytes());
                    Context.put(ip,encodedJsonStringARG1);
                    ipAndTime.put(ip, (time));
                    OrignalTimer.put(ip,time);
                    k++;
                }
                System.out.println(ipAndTime);

                Timer t = new Timer();
                t.schedule(new TimerTask() {
                    @Override
                    public void run() {
                        for (Map.Entry<String, String> mapElement : ipAndTime.entrySet()) {
                            int time = Integer.parseInt(mapElement.getValue());
                            if (time <= 0) {
                                Main.functionProcess(Context.get(mapElement.getKey()));
                                ipAndTime.put(mapElement.getKey(),OrignalTimer.get(mapElement.getKey()));
                                //System.out.println(mapElement.getValue());
                            } else {
                                time = time - 10000;
                                String t = String.valueOf(time);
                                ipAndTime.put(mapElement.getKey(), t);
                                //System.out.println(mapElement.getValue());
                                if (time <= 0) {
                                    Main.functionProcess(Context.get(mapElement.getKey()));
                                    ipAndTime.put(mapElement.getKey(), OrignalTimer.get(mapElement.getKey()));
                                }
                            }
                        }

                    }

                },0,1000);
            }
        }
}


