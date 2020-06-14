package util.Restful;

public class PivotPoint {
	final public double r2;
	final public double r1;
	final public double close;
	final public double s1;
	final public double s2;
	public Boolean aboveS1 = null;
	public Boolean aboveS2 = null;	
	public Boolean aboveR2 = null;	
	public Boolean aboveR1 = null;	
	public PivotPoint(double open, double close, double high, double low) {
		double p = (high + low + close) / 3;
		this.r1 = 2 * p - low;
		this.r2 = p + (high - low);
		this.s1 = 2 * p - high;
		this.s2 = p - (high - low);
		this.close = close;
	}
}
 