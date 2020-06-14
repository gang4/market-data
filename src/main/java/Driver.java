
import java.io.IOException;
import java.util.List;

import org.apache.commons.cli.ParseException;

import data.yahoo.entity.YahooEntity;
import data.yahoo.util.YahooDataBuilder;
import util.Restful.Download;
import util.Restful.MarketData;
import util.Restful.DataBuilder;

public class Driver {
	static int s1Count;
	static int s2Count;
	static int r1Count;
	static int r2Count;
	static int count;
	static public void main(String [] args) {
		CmdCfg cfg = new CmdCfg();
		try {
			cfg.read(args);
		} catch (ParseException e1) {
			e1.printStackTrace();
		}
		DataBuilder<YahooEntity> builder = new YahooDataBuilder(cfg.symbol, cfg.range, cfg.interval);
		Download<YahooEntity> dl = new Download<YahooEntity>(builder.getUrl());
		try {
			YahooEntity entity = dl.download(YahooEntity.class);
			if (entity.chart.error != null) {
				System.out.println(entity.chart.error);
				return;
			}

			List<MarketData> l = builder.getPivotPoints(entity);
			Driver.s1Count = 0;
			Driver.s2Count = 0;	
			Driver.r1Count = 0;
			Driver.r2Count = 0;	
			Driver.count = 0;
			System.out.println("open, lose, high, low, r2, r1, last, s1, s2, volume");
			l.forEach(d -> {
				Driver.count ++;
				//
				if (d.pp.aboveS1 != null && !d.pp.aboveS1 && d.pp.aboveS2) {
					Driver.s1Count ++;
				}
				if (d.pp.aboveS1 != null && !d.pp.aboveS2) {
					Driver.s2Count ++;
				}
				//
				if (d.pp.aboveS1 != null && d.pp.aboveR1 && !d.pp.aboveR2) {
					Driver.r1Count ++;
				}
				if (d.pp.aboveS1 != null && d.pp.aboveR2) {
					Driver.r2Count ++;
				}				
				System.out.printf("%d, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %5.2f, %d, %b, %b, %b, %b%n", d.tempstamp, d.dp.open, d.dp.close, d.dp.high, d.dp.low, d.pp.r2, d.pp.r1, d.pp.s1, d.pp.s2, d.dp.volume, d.pp.aboveS1, d.pp.aboveS2, d.pp.aboveR1, d.pp.aboveR2);
			});
			
			System.out.printf("Total: %d, s1 < low < s2: %d low < s2: %d, r1 < high < r2: %d, r2 < high %d%n", Driver.count, Driver.s1Count, Driver.s2Count, Driver.s1Count, Driver.r2Count);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
}
