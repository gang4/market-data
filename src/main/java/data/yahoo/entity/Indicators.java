package data.yahoo.entity;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"quote",
"adjclose"
})
public class Indicators {

@JsonProperty("quote")
public List<Quote> quote = null;
@JsonProperty("adjclose")
public List<Adjclose> adjclose = null;

}
