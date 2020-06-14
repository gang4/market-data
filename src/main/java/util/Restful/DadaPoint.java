package util.Restful;

public class DadaPoint {
	final public double open;
	final public double high;
	final public double low;
	final public double close;
	final public int volume;
	public DadaPoint(double open, double close, double high, double low, int volume) {
		this.open = open;
		this.high = high;
		this.low = low;
		this.close = close;
		this.volume = volume;
	}
}
