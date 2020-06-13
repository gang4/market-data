package data.yahoo.entity;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"timezone",
"start",
"end",
"gmtoffset"
})
public class Pre {

@JsonProperty("timezone")
public String timezone;
@JsonProperty("start")
public Integer start;
@JsonProperty("end")
public Integer end;
@JsonProperty("gmtoffset")
public Integer gmtoffset;

}
