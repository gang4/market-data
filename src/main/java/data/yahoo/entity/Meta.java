package data.yahoo.entity;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"currency",
"symbol",
"exchangeName",
"instrumentType",
"firstTradeDate",
"regularMarketTime",
"gmtoffset",
"timezone",
"exchangeTimezoneName",
"regularMarketPrice",
"chartPreviousClose",
"priceHint",
"currentTradingPeriod",
"dataGranularity",
"range",
"validRanges"
})
public class Meta {

@JsonProperty("currency")
public String currency;
@JsonProperty("symbol")
public String symbol;
@JsonProperty("exchangeName")
public String exchangeName;
@JsonProperty("instrumentType")
public String instrumentType;
@JsonProperty("firstTradeDate")
public Integer firstTradeDate;
@JsonProperty("regularMarketTime")
public Integer regularMarketTime;
@JsonProperty("gmtoffset")
public Integer gmtoffset;
@JsonProperty("timezone")
public String timezone;
@JsonProperty("exchangeTimezoneName")
public String exchangeTimezoneName;
@JsonProperty("regularMarketPrice")
public Double regularMarketPrice;
@JsonProperty("chartPreviousClose")
public Double chartPreviousClose;
@JsonProperty("priceHint")
public Integer priceHint;
@JsonProperty("currentTradingPeriod")
public CurrentTradingPeriod currentTradingPeriod;
@JsonProperty("dataGranularity")
public String dataGranularity;
@JsonProperty("range")
public String range;
@JsonProperty("validRanges")
public List<String> validRanges = null;

}
