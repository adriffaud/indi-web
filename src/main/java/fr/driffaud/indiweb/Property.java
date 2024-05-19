package fr.driffaud.indiweb;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class Property {
    public String device;
    public String group;
    public PropertyType type;
    public String name;
    public String label;
    public String state;
    public String perm;
    public String timeout;
    public String timestamp;
    public List<Value> values = new ArrayList<>();

    public Property() {
    }

    public Property(String device, String group, PropertyType type, String name, String label, String state, String perm, String timeout, String timestamp, List<Value> values) {
        this.device = device;
        this.group = group;
        this.type = type;
        this.name = name;
        this.label = label;
        this.state = state;
        this.perm = perm;
        this.timeout = timeout;
        this.timestamp = timestamp;
        this.values = values;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Property property)) return false;
        return Objects.equals(device, property.device) && Objects.equals(group, property.group) && type == property.type && Objects.equals(name, property.name) && Objects.equals(label, property.label) && Objects.equals(state, property.state) && Objects.equals(perm, property.perm) && Objects.equals(timeout, property.timeout) && Objects.equals(timestamp, property.timestamp) && Objects.equals(values, property.values);
    }

    @Override
    public int hashCode() {
        return Objects.hash(device, group, type, name, label, state, perm, timeout, timestamp, values);
    }

    @Override
    public String toString() {
        return "Property{" +
                "device='" + device + '\'' +
                ", group='" + group + '\'' +
                ", type=" + type +
                ", name='" + name + '\'' +
                ", label='" + label + '\'' +
                ", state='" + state + '\'' +
                ", perm='" + perm + '\'' +
                ", timeout='" + timeout + '\'' +
                ", timestamp='" + timestamp + '\'' +
                ", values=" + values +
                '}';
    }
}