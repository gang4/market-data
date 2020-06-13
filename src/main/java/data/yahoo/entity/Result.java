package data.yahoo.entity;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"meta",
"timestamp",
"indicators"
})
public class Result {

@JsonProperty("meta")
public Meta meta;
@JsonProperty("timestamp")
public List<Integer> timestamp = null;
@JsonProperty("indicators")
public Indicators indicators;

}
