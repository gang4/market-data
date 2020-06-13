package data.yahoo.entity;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"low",
"high",
"close",
"open",
"volume"
})
public class Quote {

@JsonProperty("low")
public List<Double> low = null;
@JsonProperty("high")
public List<Double> high = null;
@JsonProperty("close")
public List<Double> close = null;
@JsonProperty("open")
public List<Double> open = null;
@JsonProperty("volume")
public List<Integer> volume = null;

}
