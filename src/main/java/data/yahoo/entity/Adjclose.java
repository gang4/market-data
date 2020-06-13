package data.yahoo.entity;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"adjclose"
})
public class Adjclose {
	@JsonProperty("adjclose")
	public List<Double> adjclose = null;
}
