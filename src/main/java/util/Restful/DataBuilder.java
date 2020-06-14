package util.Restful;

import java.util.List;

public interface DataBuilder<T> {
	public String getUrl();
	public List<MarketData> getPivotPoints(T entity);
}
