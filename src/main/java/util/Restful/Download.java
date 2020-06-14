package util.Restful;

import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.net.URLConnection;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectReader;

public class Download<T> {
	final private String url;
	public Download(String _url) {
		this.url = _url;
	}

	public <T> T download(Class<T> cls) throws IOException {
		URL url = new URL(this.url);
		URLConnection conn = url.openConnection();
		InputStream is = conn.getInputStream();
		
		ObjectMapper objectMapper = new ObjectMapper();
		ObjectReader reader = objectMapper.readerFor(cls);
		T entity = reader.readValue(is);
		is.close();
		return entity;
	}
}
