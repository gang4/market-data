package util.Restful;

import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.net.URLConnection;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectReader;

public class Download<T> {
	static private final String src = "https://query1.finance.yahoo.com/v7/finance/chart/";
	// slb?range=3mo&interval=1d";
	final private String symbol;
	final private String range;
	final private String interval;	
	public Download(String symbol, String range, String interval) {
		this.symbol = symbol;
		this.range = range;
		this.interval = interval;
	}
	
	public <T> T download(Class<T> cls) throws IOException {
		URL url = new URL(src + this.symbol + "?range=" + this.range + "&interval=" + this.interval);
		URLConnection conn = url.openConnection();
		InputStream is = conn.getInputStream();
		
		ObjectMapper objectMapper = new ObjectMapper();
		ObjectReader reader = objectMapper.readerFor(cls);
		T entity = reader.readValue(is);
		is.close();
		return entity;
	}
}
