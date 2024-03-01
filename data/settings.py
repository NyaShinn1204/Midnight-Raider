import tkinter as tk

class Setting:
  language_variable = tk.StringVar()
  language_variable.set("")

  # enabled ck.button
  proxy_enabled = tk.BooleanVar()
  proxy_enabled.set(False)


  tokens = []
  validtoken = 0
  invalidtoken = 0
  lockedtoken = 0
  
  token_filenameLabel = tk.StringVar()
  token_filenameLabel.set("")
  totaltokenLabel = tk.StringVar()
  totaltokenLabel.set("Total: 000")
  validtokenLabel = tk.StringVar()
  validtokenLabel.set("Valid: 000")
  invalidtokenLabel = tk.StringVar()
  invalidtokenLabel.set("Invalid: 000")
  lockedtokenLabel = tk.StringVar()
  lockedtokenLabel.set("Locked: 000")
  
  
  proxytype = tk.StringVar()
  proxytype.set("http")
  proxies = []
  totalproxies = 0
  vaildproxies = 0
  invaildproxies = 0

  proxysetting = tk.BooleanVar()
  proxysetting.set(False)

  proxy_filenameLabel = tk.StringVar()
  proxy_filenameLabel.set("")
  totalProxiesLabel = tk.StringVar()
  totalProxiesLabel.set("Total: 000")
  validProxiesLabel = tk.StringVar()
  validProxiesLabel.set("Valid: 000")
  invalidProxiesLabel = tk.StringVar()
  invalidProxiesLabel.set("Invalid: 000")
  
  
  # joiner
  suc_joiner_Label = tk.StringVar()
  suc_joiner_Label.set("Success: 000")
  fai_joiner_Label = tk.StringVar()
  fai_joiner_Label.set("Failed: 000")
  
  joiner_link = tk.StringVar()
  joiner_link.set("")
  joiner_bypassms = tk.BooleanVar()
  joiner_bypassms.set(False)
  joiner_bypasscap = tk.BooleanVar()
  joiner_bypasscap.set(False)
  joiner_deletems = tk.BooleanVar()
  joiner_deletems.set("False")
  joiner_serverid = tk.StringVar()
  joiner_serverid.set("")
  joiner_channelid = tk.StringVar()
  joiner_channelid.set("")
  
  joiner_delay = tk.DoubleVar()
  joiner_delay.set(0.1)  
  
  # leaver
  suc_leaver_Label = tk.StringVar()
  suc_leaver_Label.set("Success: 000")
  fai_leaver_Label = tk.StringVar()
  fai_leaver_Label.set("Failed: 000")
  
  leaver_serverid = tk.StringVar()
  leaver_serverid.set("")

  leaver_delay = tk.DoubleVar()
  leaver_delay.set(0.1)  

  # normal spam
  suc_nmspam_Label = tk.StringVar()
  suc_nmspam_Label.set("Success: 000")
  fai_nmspam_Label = tk.StringVar()
  fai_nmspam_Label.set("Failed: 000")
  
  nmspam_allping = tk.BooleanVar()
  nmspam_allping.set(False)
  nmspam_allch = tk.BooleanVar()
  nmspam_allch.set(False)
  nmspam_rdstring = tk.BooleanVar()
  nmspam_rdstring.set(False)
  nmspam_ratefixer = tk.BooleanVar()
  nmspam_ratefixer.set(False)
  nmspam_randomconvert = tk.BooleanVar()
  nmspam_randomconvert.set(False)
  
  nmspam_serverid = tk.StringVar()
  nmspam_serverid.set("")
  nmspam_channelid = tk.StringVar()
  nmspam_channelid.set("")
  
  nmspam_delay = tk.DoubleVar()
  nmspam_delay.set(0.1)
  nmspam_mention = tk.DoubleVar()
  nmspam_mention.set(3)
  
  
  # go spam
  suc_gospam_Label = tk.StringVar()
  suc_gospam_Label.set("Success: 000")
  fai_gospam_Label = tk.StringVar()
  fai_gospam_Label.set("Failed: 000")  
  
  gospam_allch = tk.BooleanVar()
  gospam_allch.set(False)
  gospam_allping = tk.BooleanVar()
  gospam_allping.set(False)

  gospam_serverid = tk.StringVar()
  gospam_serverid.set("")
  gospam_channelid = tk.StringVar()
  gospam_channelid.set("")
  
  gospam_threads = tk.DoubleVar()
  gospam_threads.set(25)
  gospam_delay = tk.DoubleVar()
  gospam_delay.set(0.1)
  gospam_mention = tk.DoubleVar()
  gospam_mention.set(3)
  
  
  # threads spam
  suc_threadsspam_Label = tk.StringVar()
  suc_threadsspam_Label.set("Success: 000")
  fai_threadsspam_Label = tk.StringVar()
  fai_threadsspam_Label.set("Failed: 000")  
  
  threadsspam_serverid = tk.StringVar()
  threadsspam_serverid.set("")
  threadsspam_channelid = tk.StringVar()
  threadsspam_channelid.set("")
  
  threadsspam_threads = tk.DoubleVar()
  threadsspam_threads.set(25)
  threadsspam_delay = tk.DoubleVar()
  threadsspam_delay.set(0.1)
  threadsspam_mention = tk.DoubleVar()
  threadsspam_mention.set(3)
  
  
  # reaction spam
  suc_reactionspam_Label = tk.StringVar()
  suc_reactionspam_Label.set("Success: 000")
  fai_reactionspam_Label = tk.StringVar()
  fai_reactionspam_Label.set("Failed: 000")
  
  reaction_channelid = tk.StringVar()
  reaction_channelid.set("")
  reaction_messageid = tk.StringVar()
  reaction_messageid.set("")
  reaction_emoji = tk.StringVar()
  reaction_emoji.set("")
  
  reacspam_delay = tk.DoubleVar()
  reacspam_delay.set(0.1)
  
  
class SettingVariable:
  joinerresult_success = 0
  joinerresult_failed = 0
  leaverresult_success = 0
  leaverresult_failed = 0
  nmspamresult_success = 0
  nmspamresult_failed = 0
  gospamresult_success = 0
  gospamresult_failed = 0
  threadsspamresult_success = 0
  threadsspamresult_failed = 0
  reactionspamresult_success = 0
  reactionspamresult_failed = 0