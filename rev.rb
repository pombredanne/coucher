require 'rubygems'
require 'restclient'
require 'json'
DB="http://127.0.0.1:5984/conflict_test"

# Write multiple documents as all_or_nothing, can introduce conflicts
def writem(docs)
  JSON.parse(RestClient.post("#{DB}/_bulk_docs", {
    "all_or_nothing" => true,
    "docs" => docs,
  }.to_json))
end

# Write one document, return the rev
def write1(doc, id=nil, rev=nil)
  doc['_id'] = id if id
  doc['_rev'] = rev if rev
  writem([doc]).first['rev']
end

# Read a document, return *all* revs
def read1(id)
  retries = 0
  loop do
    # FIXME: escape id
    res = [JSON.parse(RestClient.get("#{DB}/#{id}?conflicts=true"))]
    if revs = res.first.delete('_conflicts')
      begin
        revs.each do |rev|
          res << JSON.parse(RestClient.get("#{DB}/#{id}?rev=#{rev}"))
        end
      rescue
        retries += 1
        raise if retries >= 5
        next
      end
    end
    return res
  end
end

# Create DB
RestClient.delete DB rescue nil
RestClient.put DB, {}.to_json

# Write a document
rev1 = write1({"hello"=>"xxx"},"test")
p read1("test")

# Make three conflicting versions
write1({"hello"=>"foo"},"test",rev1)
write1({"hello"=>"bar"},"test",rev1)
write1({"hello"=>"baz"},"test",rev1)

res = read1("test")
p res

# Now let's replace these three with one
res.first['hello'] = "foo+bar+baz"
res.each_with_index do |r,i|
  unless i == 0
    r.replace({'_id'=>r['_id'], '_rev'=>r['_rev'], '_deleted'=>true})
  end
end
writem(res)

p read1("test")