package data.yahoo.entity;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonPropertyOrder({
"pre",
"regular",
"post"
})
public class CurrentTradingPeriod {

@JsonProperty("pre")
public Pre pre;
@JsonProperty("regular")
public Regular regular;
@JsonProperty("post")
public Post post;

}
