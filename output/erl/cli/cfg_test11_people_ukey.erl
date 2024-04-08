-module(cfg_test11_people_ukey).

-export[get_id/1, get_id_name/1, get_id_age_sex/1].

-record(cfg, {id, name, age, sex, items, desc}).

get_id({1}) ->
	[#cfg{id=1,name="name1",age=10,sex=1,items=[1,2,3],desc="说明1"}];
get_id({2}) ->
	[#cfg{id=2,name="name2",age=11,sex=2,items=[1,2,4],desc="说明2"}];
get_id(_) ->
	[].

get_id_name({1, "name1"}) ->
	[#cfg{id=1,name="name1",age=10,sex=1,items=[1,2,3],desc="说明1"}];
get_id_name({2, "name2"}) ->
	[#cfg{id=2,name="name2",age=11,sex=2,items=[1,2,4],desc="说明2"}];
get_id_name(_) ->
	[].

get_id_age_sex({1, 10, 1}) ->
	[#cfg{id=1,name="name1",age=10,sex=1,items=[1,2,3],desc="说明1"}];
get_id_age_sex({2, 11, 2}) ->
	[#cfg{id=2,name="name2",age=11,sex=2,items=[1,2,4],desc="说明2"}];
get_id_age_sex(_) ->
	[].


