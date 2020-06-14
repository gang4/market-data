import java.io.IOException;
import java.util.List;

import data.yahoo.entity.YahooEntity;
import data.yahoo.util.YahooDataBuilder;
import util.Restful.Download;
import util.Restful.MarketData;
import util.Restful.PivotPoint;
import util.Restful.DataBuilder;

public class Driver {

	static public void main(String [] args) {
		DataBuilder<YahooEntity> builder = new YahooDataBuilder("slb", "3mo", "1d");
		Download<YahooEntity> dl = new Download<YahooEntity>(builder.getUrl());
		try {
			YahooEntity entity = dl.download(YahooEntity.class);
			List<MarketData> l = builder.getPivotPoints(entity);
			l.forEach(d -> {
				System.out.printf("%d, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %d%n", d.tempstamp, d.dp.open, d.dp.close, d.dp.high, d.dp.low, d.pp.r2, d.pp.r1, d.pp.s1, d.pp.s2, d.dp.volume);
			});
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
}
