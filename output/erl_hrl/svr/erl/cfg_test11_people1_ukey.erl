-module(cfg_test11_people1_ukey).

-export[get_id/1, get_name/1, get_age_name/1, get_id_name/1, get_id_age_sex/1].

-include("../include/cfg.hrl").

get_id({1}) ->
	[#cfg_test11_people1_ukey{id=1,name="name2",age=10,sex=1,items=[1,2,3],desc="说明1"}];
get_id({2}) ->
	[#cfg_test11_people1_ukey{id=2,name="name2",age=10,sex=2,items=[1,2,4],desc="说明2"}];
get_id(_) ->
	[].

get_name({"name2"}) ->
	[#cfg_test11_people1_ukey{id=1,name="name2",age=10,sex=1,items=[1,2,3],desc="说明1"},
	#cfg_test11_people1_ukey{id=2,name="name2",age=10,sex=2,items=[1,2,4],desc="说明2"}];
get_name(_) ->
	[].

get_age_name({10, "name2"}) ->
	[#cfg_test11_people1_ukey{id=1,name="name2",age=10,sex=1,items=[1,2,3],desc="说明1"},
	#cfg_test11_people1_ukey{id=2,name="name2",age=10,sex=2,items=[1,2,4],desc="说明2"}];
get_age_name(_) ->
	[].

get_id_name({1, "name2"}) ->
	[#cfg_test11_people1_ukey{id=1,name="name2",age=10,sex=1,items=[1,2,3],desc="说明1"}];
get_id_name({2, "name2"}) ->
	[#cfg_test11_people1_ukey{id=2,name="name2",age=10,sex=2,items=[1,2,4],desc="说明2"}];
get_id_name(_) ->
	[].

get_id_age_sex({1, 10, 1}) ->
	[#cfg_test11_people1_ukey{id=1,name="name2",age=10,sex=1,items=[1,2,3],desc="说明1"}];
get_id_age_sex({2, 10, 2}) ->
	[#cfg_test11_people1_ukey{id=2,name="name2",age=10,sex=2,items=[1,2,4],desc="说明2"}];
get_id_age_sex(_) ->
	[].


