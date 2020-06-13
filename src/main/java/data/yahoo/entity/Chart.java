package data.yahoo.entity;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"result",
"error"
})
public class Chart {

@JsonProperty("result")
public List<Result> result = null;
@JsonProperty("error")
public Object error;

}
