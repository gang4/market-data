package data.yahoo.util;

import java.util.ArrayList;
import java.util.List;

import data.yahoo.entity.Adjclose;
import data.yahoo.entity.Quote;
import data.yahoo.entity.YahooEntity;
import util.Restful.DadaPoint;
import util.Restful.DataBuilder;
import util.Restful.MarketData;
import util.Restful.PivotPoint;

public class YahooDataBuilder implements DataBuilder<YahooEntity> {
	static private final String src = "https://query1.finance.yahoo.com/v7/finance/chart/";
	// slb?range=3mo&interval=1d";
	final private String symbol;
	final private String range;
	final private String interval;	
	public YahooDataBuilder(String symbol, String range, String interval) {
		this.symbol = symbol;
		this.range = range;
		this.interval = interval;
	}
	
	public String getUrl() {
		return YahooDataBuilder.src + this.symbol + "?range=" + this.range + "&interval=" + this.interval;
	}

	public List<MarketData> getPivotPoints(YahooEntity entity) {
		int len = entity.chart.result.get(0).timestamp.size();
		final List<MarketData> list = new ArrayList<>();
		Quote q = entity.chart.result.get(0).indicators.quote.get(0);
		Adjclose c = entity.chart.result.get(0).indicators.adjclose.get(0);
		for (int i = 0; i < len; i ++) {
			MarketData md = new MarketData();
			md.pp = new PivotPoint(q.open.get(i), c.adjclose.get(i), q.high.get(i), q.low.get(i));
			md.dp = new DadaPoint(q.open.get(i), c.adjclose.get(i), q.high.get(i), q.low.get(i), q.volume.get(i));
			md.tempstamp = entity.chart.result.get(0).timestamp.get(i);
			list.add(md);
		}
		return list;
	}
}
