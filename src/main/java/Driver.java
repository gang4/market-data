import java.io.IOException;

import data.yahoo.entity.YahooEntity;
import util.Restful.Download;

public class Driver {

	static public void main(String [] args) {
		Download<YahooEntity> dl = new Download<YahooEntity>("slb", "3mo", "1d");
		try {
			YahooEntity entity = dl.download(YahooEntity.class);
			System.out.println(entity);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
}
